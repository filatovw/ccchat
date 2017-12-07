Test task "chat"
================

## Task

Write go client and server for chatting, that communicates over web socket.

Client should establish connection to the server and send text messages. Messages format is:

    command_code[::[msg]]\r\n 

Possible values for `command_code`:

    auth - authorization, msg - client name,
    end - close session, msg and `::` should be omitted
    some_key - arbitrary key (may be empty), msg - arbitrary message.

msg - argitrary text, can be empty

1. Server should receive messages and render them in format 
`[client_name]: some_key | msg`.
`client_name` - name of a client, that is placed in `msg` of authorization message, `some_key` - key of message server recieves after authorization, `msg` - message for corresponding key.
2. You need to create a page (web interface) with the following elements: list of messages, that updates in real time, list of clients that updates in real time, number of connected clients, current number of messages.
3. All given messages should be stored in your favourite database and should be displayed on connection to the web interface.
4. Client should have configuration mechanism for automatic messages generation (number of messages per session, session duration, any other options depending on your choice). Configuration should be stored in a file of `JSON` format or should be passed as arguments in command line.

## Example

Client sends:

    auth::Jackie Chan
    k1::Hello!
    k2::My name is Jackie Chan
    k3::What is your name?
    end
    k4::Bye!

Result displayed by server:

Messages: 2 | Clients: 1
------------|-----------
[Jackie Chan]: k1 \| Hello! | Empty
[Jackie Chan]: k2 \| My name is Jackie Chan |
[Jachie Chan]: k3 \| What is your name? |