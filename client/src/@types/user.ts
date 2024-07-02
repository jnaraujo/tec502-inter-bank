export interface User {
  id: number
  name: string
  documents: string[]
  balance: string
  createdAt: Date
  ibk: string
  type: "individual" | "legal" | "joint"
}
