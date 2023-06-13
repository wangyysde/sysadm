running mode:
	command: agent run as CLI  with commands and flags.agent send result messages to a server,stdout or  file when it running.
	daemon: agent running as daemon process after start. agent gets or receives commands (depending on get mode),then run then on the host and reponses 
	        the result to a server 
		getmode:
		    passive: agent gets commands from the server periodically and run them.
			active: the server send a command to a agent when it  want the agent to run the command on a host. 