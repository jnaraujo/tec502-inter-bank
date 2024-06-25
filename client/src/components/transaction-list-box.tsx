import { TransactionBox } from "./transaction-box"

export function TransactionListBox() {
  return (
    <article className="flex flex-col justify-between space-y-4 overflow-auto rounded-lg border bg-background p-6">
      <h2 className="text-xl font-medium text-zinc-900">Suas transações</h2>

      <div className="flex flex-col gap-6 overflow-auto py-2">
        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="pending"
          operations={[
            {
              from: "1-1",
              to: "2-1",
              amount: 200,
            },
            {
              from: "1-1",
              to: "2-2",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="success"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="failed"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="pending"
          operations={[
            {
              from: "1-1",
              to: "2-1",
              amount: 200,
            },
            {
              from: "1-1",
              to: "2-2",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="success"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="failed"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="pending"
          operations={[
            {
              from: "1-1",
              to: "2-1",
              amount: 200,
            },
            {
              from: "1-1",
              to: "2-2",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="success"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="failed"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="pending"
          operations={[
            {
              from: "1-1",
              to: "2-1",
              amount: 200,
            },
            {
              from: "1-1",
              to: "2-2",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="success"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />

        <TransactionBox
          createdAt="2024-06-25T16:35:57.794293059-03:00"
          status="failed"
          operations={[
            {
              from: "2-1",
              to: "1-1",
              amount: 250,
            },
          ]}
        />
      </div>
    </article>
  )
}
