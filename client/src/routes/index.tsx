import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Welcome to gRPC CRUD App
        </h1>
        <p className="text-lg text-gray-600">
          Built with gRPC, Connect, Go, Vite, TanStack Router, and TailwindCSS
        </p>
      </div>
    </div>
  )
}
