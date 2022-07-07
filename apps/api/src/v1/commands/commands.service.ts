import { HttpException, Injectable } from '@nestjs/common';
import { Client, Transport } from '@nestjs/microservices';
import { config } from '@tsuwari/config';
import { Command, PrismaService, Response } from '@tsuwari/prisma';
import { ClientProxy } from '@tsuwari/shared';

import { RedisService } from '../../redis.service.js';
import { UpdateOrCreateCommandDto } from './dto/create.js';

@Injectable()
export class CommandsService {
  @Client({ transport: Transport.NATS, options: { servers: [config.NATS_URL] } })
  nats: ClientProxy;

  constructor(
    private readonly prisma: PrismaService,
    private readonly redis: RedisService,
  ) { }


  async getList(userId: string) {
    const commands: (Command & {
      responses?: Response[];
    })[] = await this.prisma.command.findMany({
      where: { channelId: userId },
      include: {
        responses: true,
      },
    });

    const defaultCommands = await this.nats.send('bots.getDefaultCommands', {}).toPromise();
    if (defaultCommands) {
      for (const command of defaultCommands) {
        if (!commands.some(c => c.defaultName === command.name)) {
          const newCommand = await this.prisma.command.create({
            data: {
              channelId: userId,
              default: true,
              defaultName: command.name,
              description: command.description,
              visible: command.visible,
              name: command.name,
              permission: command.permission,
              cooldown: 0,
              cooldownType: 'GLOBAL',
            },
          });

          this.setCommandCache(newCommand);
          commands.push(newCommand);
        }
      }
    }

    console.log(defaultCommands, commands);

    return commands;
  }

  async setCommandCache(command: Command & { responses?: Response[] }, oldCommand?: Command & { responses?: Response[] }) {
    const preKey = `commands:${command.channelId}`;

    if (oldCommand) {
      await this.redis.del(`commands:${oldCommand.channelId}:${oldCommand.name}`);

      if (oldCommand.aliases && Array.isArray(command.aliases)) {
        for (const alias of command.aliases) {
          await this.redis.del(`${preKey}:${alias}`);
        }
      }
    }

    const commandForSet = {
      ...command,
      responses: command.responses ? JSON.stringify(command.responses.map(r => r.text)) : [],
      aliases: Array.isArray(command.aliases) ? JSON.stringify(command.aliases) : command.aliases,
    };

    await this.redis.hmset(`${preKey}:${command.name}`, commandForSet);

    if (command.aliases && Array.isArray(command.aliases)) {
      for (const alias of command.aliases) {
        await this.redis.hmset(`${preKey}:${alias}`, commandForSet);
      }
    }

  }

  async create(userId: string, data: UpdateOrCreateCommandDto & { defaultName?: string }) {
    const isExists = await this.prisma.command.findMany({
      where: {
        name: data.name,
        OR: {
          name: { in: data.aliases },
          aliases: {
            array_contains: data.aliases,
          },
        },
      },
    });

    if (isExists.length) {
      throw new HttpException(`Command already exists`, 400);
    }

    if (!data.responses?.length) {
      throw new HttpException(`You should add atleast 1 response to command.`, 400);
    }

    const command = await this.prisma.command.create({
      data: {
        ...data,
        channelId: userId,
        responses: {
          createMany: {
            data: data.responses.filter(r => r.text).map((r) => ({ text: r.text?.trim().replace(/(\r\n|\n|\r)/, '') })),
          },
        },
      },
      include: {
        responses: true,
      },
    });

    await this.setCommandCache(command);
    return command;
  }

  async delete(userId: string, commandId: string) {
    const command = await this.prisma.command.findFirst({ where: { channelId: userId, id: commandId } });

    if (!command) {
      throw new HttpException('Command not exists', 404);
    }

    if (command.default) {
      throw new HttpException('You cannot delete default command.', 400);
    }

    const result = await this.prisma.command.delete({
      where: {
        id: commandId,
      },
    });

    await this.redis.del(`commands:${userId}:${command.name}`);
    if (Array.isArray(command.aliases)) {
      for (const aliase of command.aliases as string[]) {
        await this.redis.del(`commands:${userId}:${aliase}`);
      }
    }

    return result;
  }

  async update(userId: string, commandId: string, data: UpdateOrCreateCommandDto) {
    const command = await this.prisma.command.findFirst({
      where: { channelId: userId, id: commandId },
      include: { responses: true },
    });

    if (!command) {
      throw new HttpException('Command not exists', 404);
    }

    if (!command.responses?.length && !command.default) {
      throw new HttpException(`You should add atleast 1 response to command.`, 400);
    }

    data.responses = data.responses?.filter(r => r.text).map(r => ({ ...r, text: r.text ? r.text.trim().replace(/(\r\n|\n|\r)/, '') : null }));

    const responsesForUpdate = data.responses
      .filter(r => command.responses.some(c => c.id === r.id && r.text && r.id))
      .map(r => ({ id: r.id, text: r.text }))
      .map(r => this.prisma.response.update({ where: { id: r.id }, data: { text: r.text } }));

    const [newCommand] = await this.prisma.$transaction([
      this.prisma.command.update({
        where: { id: commandId },
        data: {
          ...data,
          channelId: userId,
          responses: {
            deleteMany: command.responses.filter(r => !data.responses.map(s => s.id).includes(r.id)),
            createMany: {
              data: data.responses
                .filter((r) => !command.responses.some((c) => c.id === r.id)),
              skipDuplicates: true,
            },
          },
        },
        include: {
          responses: true,
        },
      }),
      ...responsesForUpdate,
    ]);

    const newResponses = await this.prisma.response.findMany({ where: { commandId: command.id } });

    await this.setCommandCache({
      ...newCommand,
      responses: newResponses.flat(),
    }, command);

    return {
      ...newCommand,
      responses: newResponses.flat(),
    };
  }
}
