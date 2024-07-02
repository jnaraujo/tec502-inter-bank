export interface Operation {
  id: string
  from: string
  to: string
  type: string
  amount: string
  status: string
  created_at: string
  updated_at: string
}

export interface Transaction {
  id: string
  parent_id: string
  author: string
  operations: Operation[]
  created_at: string
  updated_at: string
  status: string
  type: "package" | "final"
}
