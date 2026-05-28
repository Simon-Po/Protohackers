opts = [
      :binary,
      packet: 0,
      active: false,
    ]


{:ok,sock} = :gen_tcp.connect(String.to_charlist("178.105.250.71"), 8080, opts)

# Byte:  |  0  |  1     2     3     4  |  5     6     7     8  |
# Type:  |char |         int32         |         int32         |
# Value: | 'I' |       timestamp       |         price         |
msg = <<?I ,12345::signed-32-big,101::signed-32-big>>
msg2 = <<?I ,12346::signed-32-big,102::signed-32-big>>
msg3 = <<?I ,12347::signed-32-big,100::signed-32-big>>
msg4 = <<?I ,40960::signed-32-big,5::signed-32-big>>
msg5 = <<?Q ,12288::signed-32-big,16384::signed-32-big>>

:ok = :gen_tcp.send(sock,msg)
:ok = :gen_tcp.send(sock,msg2)
:ok = :gen_tcp.send(sock,msg3)
:ok = :gen_tcp.send(sock,msg4)
:ok = :gen_tcp.send(sock,msg5)

case :gen_tcp.recv(sock, 4, 5_000) do
  {:ok, data} ->
    IO.inspect(data, label: "response")
    <<value::signed-big-32>> = data
    expected = 101

    if value == expected do
      IO.puts("PASS: got #{value}")
    else
      IO.puts("FAIL: expected #{expected}, got #{value}")
    end

  {:error, :timeout} ->
    IO.puts("no response after 5 seconds")
  {:error, :closed} ->
    IO.puts("server closed connection")
end
:gen_tcp.close(sock)
