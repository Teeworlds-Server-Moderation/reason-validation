# reason-validation

This is a microservice that is part of the Teeworlds Server Moderation Framework.
This one is supposed to receive events from the broker and act according to a predefined configuration and classification of voting reasons and their corresponding actions/reactions in a csv file.

If someone on any moderated servers does start a specvote, trying to move another player to the spectators, or starts a kickvote, trying to kick a player from a server, this microservice gets both of the voting events and evaluates the reasoning. If the reason matches any reasons in the classicifation file, the microservice acct accordingly, either ignores, votebans or aborts the vote.

This will be a basic proof of concept to see if this actually helps in decreasing the numbe rof funvotes on the servers.
An interesting extension would be some machine learning attempt that I cannot really implement, as I do not have enough knowldge on that topic.

Anayway, I did classify about 4500 unique reason strings that were logged in the past 1.5 years on my zCatch servers and will see how far this journey goes.
