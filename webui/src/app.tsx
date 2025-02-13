"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"

interface Container {
  ContainerID: string
  ContainerName: string
  ContainerIP: string
}

interface NetworkData {
  HallucinetNetwork: string
  Networks: {
    [key: string]: Container[]
  }
}

export default function NetworkTables() {
  const [data, setData] = useState<NetworkData | null>(null)
  const [error, setError] = useState<string>("")

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080")
        if (!response.ok) throw new Error("Network response was not ok")
        const networkData = await response.json()
        setData(networkData)
      } catch (err) {
        setError("Failed to fetch network data")
        console.error(err)
      }
    }

    fetchData()
  }, [])

  if (error) return <div className="text-red-500 p-4">{error}</div>
  if (!data) return <div className="p-4">Loading...</div>

  const hallucinetContainers = data.Networks[data.HallucinetNetwork] || []
  const otherNetworks = Object.entries(data.Networks).filter(([name]) => name !== data.HallucinetNetwork)

  return (
    <div className="p-4 space-y-8 bg-black text-white min-h-screen">
      <section>
        <h1 className="text-4xl font-bold text-center mb-4">Hallucinet</h1>
        <p className="text-center text-gray-400 mb-4">
          Containers connected to the network{" "}
          <span className="bg-gray-800 px-2 py-1 rounded">{data.HallucinetNetwork}</span>
        </p>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="text-gray-400">Domain</TableHead>
              <TableHead className="text-gray-400">Hallucinet IP</TableHead>
              <TableHead className="text-gray-400">ID</TableHead>
              <TableHead className="text-gray-400">Name</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {hallucinetContainers.map((container) => (
              <TableRow key={container.ContainerID}>
                <TableCell>.test</TableCell>
                <TableCell>{container.ContainerIP}</TableCell>
                <TableCell className="font-mono">{container.ContainerID.slice(0, 12)}</TableCell>
                <TableCell>{container.ContainerName}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </section>

      <section>
        <h2 className="text-4xl font-bold text-center mb-6">Other Networks</h2>
        <div className="space-y-6">
          {otherNetworks.map(([networkName, containers]) => (
            <div key={networkName} className="bg-gray-900 rounded-lg p-4">
              <h3 className="text-xl font-semibold mb-4 px-4">{networkName}</h3>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="text-gray-400">IP</TableHead>
                    <TableHead className="text-gray-400">ID</TableHead>
                    <TableHead className="text-gray-400">Name</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {containers.map((container) => (
                    <TableRow key={container.ContainerID}>
                      <TableCell>{container.ContainerIP}</TableCell>
                      <TableCell className="font-mono">{container.ContainerID.slice(0, 12)}</TableCell>
                      <TableCell>{container.ContainerName}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}


