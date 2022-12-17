/* eslint-disable import/no-cycle */
import {
  Column,
  Entity,
  Index,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
} from 'typeorm';

import { Channel } from './Channel';

export enum SettingsType {
  links = 'links',
  blacklists = 'blacklists',
  symbols = 'symbols',
  longMessage = 'longMessage',
  caps = 'caps',
  emotes = 'emotes',
}

@Index('channels_moderation_settings_channelId_type_key', ['channelId', 'type'], { unique: true })
@Entity('channels_moderation_settings', { schema: 'public' })
export class ChannelModerationSetting {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('enum', {
    name: 'type',
    enum: SettingsType,
  })
  type: SettingsType;

  @Column('boolean', { name: 'enabled', default: false })
  enabled: boolean;

  @Column('boolean', { name: 'subscribers', default: false })
  subscribers: boolean;

  @Column('boolean', { name: 'vips', default: false })
  vips: boolean;

  @Column('integer', { name: 'banTime', default: 600 })
  banTime: number;

  @Column('text', { name: 'banMessage', nullable: false, default: '' })
  banMessage: string;

  @Column('text', { name: 'warningMessage', nullable: false, default: '' })
  warningMessage: string;

  @Column('boolean', {
    name: 'checkClips',
    nullable: true,
    default: false,
  })
  checkClips: boolean | null;

  @Column('integer', {
    name: 'triggerLength',
    nullable: true,
    default: 300,
  })
  triggerLength: number | null;

  @Column('integer', {
    name: 'maxPercentage',
    nullable: true,
    default: 50,
  })
  maxPercentage: number | null;

  @Column('text', { name: 'blackListSentences', array: true, nullable: true, default: [] })
  blackListSentences: string[] | null;

  @ManyToOne(() => Channel, _ => _.moderationSettings, {
    onDelete: 'RESTRICT',
    onUpdate: 'CASCADE',
  })
  @JoinColumn([{ name: 'channelId', referencedColumnName: 'id' }])
  channel?: Channel;

  @Column('text', { name: 'channelId' })
  channelId: string;
}
