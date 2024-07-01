import { create } from "zustand"

type Store = {
  address: string
  setAddress: (addr: string) => void
}

export const useBank = create<Store>()((set) => ({
  address: "",
  setAddress: (addr: string) => {
    set(() => ({
      address: addr,
    }))
  },
}))
