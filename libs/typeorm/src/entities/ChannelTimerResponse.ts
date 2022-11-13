import {
  Column,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn,
  type Relation,
} from 'typeorm';

import { type ChannelTimer } from './ChannelTimer.js';

@Entity('channels_timers_responses', { schema: 'public' })
export class ChannelTimerResponse {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column('text')
  text: string;

  @Column('bool', { default: true })
  isAnnounce: boolean;

  @ManyToOne('ChannelTimer', 'responses', {
    onDelete: 'CASCADE',
  })
  @JoinColumn({ name: 'timerId' })
  timer?: Relation<ChannelTimer>;

  @Column()
  timerId: string;
}
