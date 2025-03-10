import { PropsWithChildren, useEffect, useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./components/ui/table";
import { Card, CardContent, CardHeader, CardTitle } from "./components/ui/card";

export function SectionHeading(props: PropsWithChildren) {
  return (
    <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl text-center m-4">
      {props.children}
    </h1>
  );
}

export function HallucinetNetwork(props: PropsWithChildren) {
  return (
    <code className="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
      {props.children}
    </code>
  );
}

interface Container {
  ContainerID: string;
  ContainerName: string;
  ContainerIP: string;
}

interface NetworkData {
  HallucinetNetwork: string;
  Networks: {
    [key: string]: Container[];
  };
}
export default function App() {
  const [data, setData] = useState<NetworkData | null>(null);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    const fetchData = async () => {
      try {
        console.log("retrieving contianers");
        const response = await fetch(
          `http://${window.location.hostname}/api/v1/containers`,
        );
        if (!response.ok) throw new Error("Network response was not ok");
        const networkData = await response.json();
        setData(networkData);
      } catch (err) {
        setError("Failed to fetch network data");
        console.error(err);
      }
    };

    fetchData();
  }, []);

  if (error) return <div className="text-red-500 p-4">{error}</div>;
  if (!data) return <></>;

  const hallucinetContainers = data.Networks[data.HallucinetNetwork] || [];
  const otherNetworks = Object.entries(data.Networks).filter(
    ([name]) => name !== data.HallucinetNetwork,
  );
  return (
    <div className="w-screen">
      <div className="max-w-7xl mx-auto p-10 min-h-screen flex flex-col">
        <p className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl text-center m-4">
          Hallucinet
        </p>
        <Card className="m-5">
          <CardContent className="mt-5">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Domain</TableHead>
                  <TableHead>IP</TableHead>
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {hallucinetContainers.map((container) => (
                  <TableRow key={container.ContainerID}>
                    <TableCell>{container.ContainerName}.test</TableCell>
                    <TableCell>{container.ContainerIP}</TableCell>
                    <TableCell className="font-mono">
                      {container.ContainerID.slice(0, 12)}
                    </TableCell>
                    <TableCell>{container.ContainerName}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
        <p className="mt-10 scroll-m-20 text-3xl font-extrabold tracking-tight lg:text-4xl text-center m-4">
          Other Networks
        </p>
        {otherNetworks.map(([networkName, containers]) => (
          <div>
            <Card className="m-5">
              <CardHeader>
                <CardTitle className="px-1">{networkName}</CardTitle>
              </CardHeader>
              <CardContent>
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
                        <TableCell className="font-mono">
                          {container.ContainerID.slice(0, 12)}
                        </TableCell>
                        <TableCell>{container.ContainerName}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>
          </div>
        ))}
      </div>
    </div>
  );
}
