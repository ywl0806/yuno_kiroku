import { z } from 'zod'
import { LOGIN_SCHEMA } from '../consts'

export type LoginForm = z.infer<typeof LOGIN_SCHEMA>
