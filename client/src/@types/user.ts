export interface User {
  id: number
  name: string
  document: string
  balance: string
  createdAt: Date
  ibk: string
  type: "individual" | "legal" | "joint"
}
