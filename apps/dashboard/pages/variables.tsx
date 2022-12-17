import { Badge, Button, Table } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { Dashboard } from '@tsuwari/shared';
import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { useState } from 'react';
import useSWR from 'swr';

import { VariableDrawer } from '../components/variables/drawer';
import { swrFetcher } from '../services/swrFetcher';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableVariable, setEditableVariable] = useState<ChannelCustomvar>({} as any);
  const [selectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  const { data: variables } = useSWR<ChannelCustomvar[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/variables` : null,
    swrFetcher,
  );

  return (
    <div>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {variables &&
            variables.map((element, idx) => (
              <tr key={element.id}>
                <td>
                  <Badge>{element.name}</Badge>
                </td>
                <td>
                  <Badge color="cyan">{element.type}</Badge>
                </td>
                <td>
                  <Button
                    onClick={() => {
                      setEditableVariable(variables[idx] as any);
                      setEditDrawerOpened(true);
                    }}
                  >
                    Edit
                  </Button>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>

      <VariableDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        variable={editableVariable}
      />
    </div>
  );
}
