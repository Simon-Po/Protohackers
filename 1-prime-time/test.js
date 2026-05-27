const net = require("net");
const HOST = "localhost";
const PORT = 8080;

const tests = [
  {
    name: "valid multiple requests",
    requests: [
      { method: "isPrime", number: 2 },
      { method: "isPrime", number: 4 },
      { method: "isPrime", number: 17 },
      { method: "isPrime", number: 7.5 },
      { method: "isPrime", number: 19, extra: true },
    ],
  },
  {
    name: "malformed: number is string",
    requests: [{ method: "isPrime", number: "123" }],
  },
  {
    name: "malformed: wrong method",
    requests: [{ method: "wrong", number: 123 }],
  },
  {
    name: "malformed: missing number",
    requests: [{ method: "isPrime" }],
  },
  {
    name: "malformed: invalid json",
    raw: ["not json"],
  },
];

function writeRequest(socket, request) {
  const line = typeof request === "string" ? request : JSON.stringify(request);

  console.log("SEND:", line);
  socket.write(line + "\n");
}

function runTest(test) {
  return new Promise((resolve) => {
    const socket = net.createConnection({ host: HOST, port: PORT }, () => {
      console.log(`\n=== ${test.name} ===`);

      const requests = test.raw || test.requests;
      requests.forEach((request) => writeRequest(socket, request));

      setTimeout(() => socket.end(), 500);
    });

    let buffer = "";

    socket.on("data", (chunk) => {
      buffer += chunk.toString("utf8");

      const lines = buffer.split("\n");
      buffer = lines.pop();

      for (const line of lines) {
        console.log("RECV:", line);
      }
    });

    socket.on("error", (err) => {
      console.error("ERROR:", err.message);
      resolve();
    });

    socket.on("close", resolve);
  });
}

async function main() {
  for (const test of tests) {
    await runTest(test);
  }
}

main();
