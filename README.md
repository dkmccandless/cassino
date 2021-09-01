# cassino

Module cassino implements the game of Cassino.

## Rules

Cassino is a fishing card game for two players using a standard 52-card deck. The goal of the game is to score points by capturing cards from the table individually according to their rank or in groups according to their total value.

### Play

Each round begins with four cards dealt to each player. In addition, the first round begins with four cards dealt face up to the table. Starting with the player opposite the dealer, players alternate turns playing one card from their hand to perform one of the following actions:

* **Trailing**: A card may be placed on the table.
* **Capturing**: A card may be used to capture any cards on the table of the same rank. In addition, a number card may be used to capture any builds of the same value, as well as any combinations of number cards and/or simple builds whose total equals the capturing card's value. Captured cards are set aside.
* **Building**: A number card may be played onto one or more other number cards on the table to create a "build" that can only be captured by a card of a specific value; see below.

When the players have played all of the cards in their hands, the dealer deals a new hand of four cards to each player to begin the next round, and so on until the deck is exhausted. At the end of the last round, any cards left on the table are awarded to the player who last captured cards.

#### Building

Building is the distinctive gameplay mechanic of Cassino. A player may create a build in two ways. Each involves playing a card onto one or more cards (or existing builds) on the table to make an indivisible pile with a specific value, which then can only be captured by a card of the same value. Face cards have no numerical value and may not be used in any builds.

A *simple* build is created by playing a card onto a combination of cards on the table such that their total equals the value of a card in the builder's hand. A simple build behaves identically to a single card of the corresponding value, and so the value of a simple build can increase if a player plays another card onto it (subject to the restriction that they must hold a card of the corresponding value in hand).

Examples:
* With a 2 on the table, a player with a 5 and a 7 in hand plays the 5 onto the 2 to "build 7".
* With an ace and a 6 on the table, a player with a 2 and a 9 in hand groups the ace and 6 together and plays the 2 on top to "build 9".
* With a simple build of 2 and 3 (building 5) on the table, a player with a 3 and an 8 in hand plays the 3 onto it to change its value to 8.
* With a simple build of two aces (building 2) and a 5 on the table, a player with a 3 and a 10 in hand adds the 5 to the build and then plays the 3 onto it to change its value to 10.

In contrast, a *compound* build contains two or more sets of cards, each of which sum to the build's value. A compound build's value is fixed. Therefore, a compound build may not be used as part of a simple build.

Examples:
* With a 7 on the table, a player with two 7s in hand plays one of them onto the 7 on the table to "build 7s".
* With a 3 and a 6 on the table, a player with two 9s in hand groups the 3 and 6 and plays a 9 onto them to "build 9s".
* With an 8 and a 10 on the table, a player with a 2 and a 10 in hand plays the 2 onto the 8 and then adds the 10 on top to "build 10s".
* With an ace, a 2, and a 5 on the table, a player with a 4 and a 6 in hand plays the 4 onto the 2 and then combines the ace and 5 on top to "build 6s".

A player may only create or modify a build if they hold a card in their hand that can capture it. Either player may capture any build, but a build is said to be *controlled* by the player who last modified it. A player who controls a build on the table may not trail.

### Scoring

Points are awarded for captured cards as follows:

* Most cards: 3 points
* Most spades: 1 point
* The 10 of diamonds ("Big Cassino"): 2 points
* The 2 of spades ("Little Cassino"): 1 point
* Each ace: 1 point

In addition, if a player captures all cards on the table (a "sweep"), they immediately score 1 point.
