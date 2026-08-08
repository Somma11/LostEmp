import { sqliteTable, text, integer } from 'drizzle-orm/sqlite-core'

export const equipamentos = sqliteTable('equipamentos', {
  id: integer('id').primaryKey(),
  marca: text('marca').notNull(),
  modelo: text('modelo').notNull(),
  status: text('status').notNull(),
  manutencao: integer('manutencao', { mode: 'boolean' }),
})