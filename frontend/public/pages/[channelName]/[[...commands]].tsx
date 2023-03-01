import { Badge, Flex, Grid, Table, Text, Tooltip } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { NextPage } from 'next';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

import { useUsersByNames } from '@/services/users';

type Command = {
  id: string
  name: string
  responses: string[]
  permissions: string[]
  cooldown: number
  cooldownType: string
  aliases: string[]
  description: null | string
}


const Commands: NextPage = () => {
  const router = useRouter();
  const { data: users } = useUsersByNames([router.query.channelName as string]);

  const {
    data: commands,
  } = useQuery({
    queryKey: ['commands', users?.at(0)?.id],
    queryFn: async (): Promise<Command[]> => {
      const req = await fetch(`/api/v1/p/commands/${users?.at(0)?.id}`);

      return req.json();
    },
    enabled: !!users?.at(0)?.id,
  });

  return (<Table highlightOnHover>
    <thead>
    <tr>
      <th>Name</th>
      <th>Description</th>
      <th>Permissions</th>
      <th>Cooldown</th>
    </tr>
    </thead>
    <tbody>
    {commands?.map((c, commandIndex) => <tr key={commandIndex}>
      <td style={{
        whiteSpace: 'nowrap',
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        maxWidth: 150,
      }}>
        <Tooltip label={[c?.name, ...c.aliases || []].join(', ')}>
          <Text truncate>
           {[c?.name, ...c.aliases || []].join(', ')}
          </Text>
        </Tooltip>
      </td>
      <td>{c.description ? c.description : c?.responses?.map((r, responseIndex) => <Text
        key={responseIndex}
        title={r}
        lineClamp={1}
        style={{ textOverflow: 'ellipsis', overflow: 'hidden' }}
      >
        {r}
      </Text>)}</td>
      <td>
        <Flex direction={'column'} gap={'xs'}>
          {c?.permissions?.map((p, i) => <Badge key={i}>{p}</Badge>)}
          {!c.permissions?.length && <Badge color={'green'}>Everyone</Badge>}
        </Flex>
      </td>
      <td>{c?.cooldown} ({c?.cooldownType?.toLowerCase().replace('_', ' ')})</td>
    </tr>)}
    </tbody>
  </Table>);
};

export default Commands;