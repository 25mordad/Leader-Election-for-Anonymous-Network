# Echo algorithm with extinction for anonymous networks
- Leader Election Algorithm for Anonymous Network
The resulting election algorithm progresses in rounds; at the start of every round, the active processes randomly select an ID and run the echo algorithm with extinction. 
Again, round numbers are used to recognize messages from earlier rounds. When a process is hit by a wave with a higher round number than its current wave, or with the same round number but a higher ID, the process becomes passive (if it was not already) and moves to that other wave.

Initially, initiators are active at the start of round Let process p be active. At the start of an election round n ≥ 0, p randomly selects an ID id p ∈ {1,...,N}, and starts a wave, tagged with n and id p. As a third parameter it adds a 0, meaning that this is not a message to its parent. This third parameter is used to report the subtree size in messages to a parent.
A process p, which is in a wave from round n with ID i, waits for a wave message to arrive, tagged with a round number n and an ID j. When this happens, p acts as follows, depending on the parameter values in this message.
• If n > n, or n = n and j>i, then p makes the sender its parent, changes to the wave in round n with ID j, and treats the message accordingly.

• If n < n, or n = n and j<i, then p purges the message.

• If n = n and j = i, then p treats the message according to the echo algorithm.

- Distributed Algorithms by Wan Fokkink-Page91

# configuration file:
- First line is the name of root for example: 127.0.0.1:8082, it means the ip address is 127.0.0.1 and the port is 8082
- After first line you should write the neighbor in each line with ip and port for example if the node has two neighbors, we write down:
127.0.0.1:8083
127.0.0.1:8084
- If the node is initiator, we should write: initiator:[nodeIP]:[nodePort], for example: initiator:127.0.0.1:8082
and the last line "size:n" shows the number of nodes in the network

# node
How can I add a node?
- Just copy the node folder and rename it to what you want. Then change the configuration file as you need.
- All node files (node.go) are the same.

#Graph for Sample:
	81 --------- 82 ---------- 83 -------- 84 ----------- 86
							   |
							   |
							   |
							   |
							   85
			


