let socket;
let userId, roomId, username;

document.getElementById("userId").value = Math.floor(Math.random() * 1000);
document.getElementById("username").value =
  "Player" + Math.floor(Math.random() * 1000);
document.getElementById("connectBtn").addEventListener("click", function () {
  userId = document.getElementById("userId").value;
  roomId = document.getElementById("roomId").value;
  username = document.getElementById("username").value;

  const wsUrl = `ws://localhost:8000/api/v1/ws/join/${roomId}?userId=${userId}&username=${username}`;
  socket = new WebSocket(wsUrl);

  socket.onopen = function () {
    document.getElementById("connectionStatus").innerText = "Connected!";
    document.getElementById("startGameBtn").disabled = false;
    document.getElementById("dealCardsBtn").disabled = true;
    document.getElementById("changeTeamBtn").disabled = false; // Enable change team button
  };

  socket.onmessage = function (event) {
    const response = JSON.parse(event.data);
    document.getElementById("response").innerText = JSON.stringify(
      response,
      null,
      2,
    );

    if (response.action === "deal_cards") {
      displayCards(
        response.payload.cards,
        response.payload.table_cards,
        response.payload.next_turn_id,
      );
    } else if (response.action === "game_started") {
      document.getElementById("dealCardsBtn").disabled =
        userId !== response.payload.dealer_id;
    } else if (response.action === "game_ended") {
      document.getElementById("dealCardsBtn").disabled = true;
    } else if (response.action === "round_over") {
      document.getElementById("dealCardsBtn").disabled =
        userId !== response.payload.dealer_id;
    } else if (response.action === "card_played") {
      displayCards(
        response.payload.player_hand,
        response.payload.table_cards,
        response.payload.next_turn_id,
      );
    } else if (response.action === "team_changed") {
      updateTeams(response.payload.teams); // Update teams on team change
    }
  };

  socket.onerror = function (error) {
    document.getElementById("connectionStatus").innerText =
      "Error: " + error.message;
  };
});

// Update function to display teams without individual buttons
function updateTeams(teams) {
  const teamListDiv = document.getElementById("teamList");
  teamListDiv.innerHTML = ""; // Clear current team display

  teams.forEach((team) => {
    const teamDiv = document.createElement("div");
    teamDiv.className = "team";
    teamDiv.innerHTML = `<strong>Team ${team.team_id}</strong> - Score: ${team.score}`;

    if (team.players) {
      const players = team.players
        .map((player) => `${player.username} (Score: ${player.score})`)
        .join(", ");
      teamDiv.innerHTML += `<div>Players: ${players}</div>`;
    } else {
      teamDiv.innerHTML += "<div>Players: None</div>";
    }

    teamListDiv.appendChild(teamDiv);
  });
}

document.getElementById("startGameBtn").addEventListener("click", function () {
  const startGameMessage = {
    action: "start_game",
    payload: { card_id: 31 },
    user_id: userId,
    room_id: roomId,
  };
  socket.send(JSON.stringify(startGameMessage));
});

document.getElementById("dealCardsBtn").addEventListener("click", function () {
  const dealCardsMessage = {
    action: "deal_cards",
    payload: { card_id: 49 },
    user_id: userId,
    room_id: roomId,
  };
  socket.send(JSON.stringify(dealCardsMessage));
});

// Event listener for the change team button
document.getElementById("changeTeamBtn").addEventListener("click", function () {
  const changeTeamMessage = {
    action: "change_team", // Keep the same action name
    payload: { team_id: null }, // Use the payload structure as required
    user_id: userId,
    room_id: roomId,
  };
  socket.send(JSON.stringify(changeTeamMessage));
});

function displayCards(playerHand, tableCards, nextTurnId) {
  const playerHandDiv = document.getElementById("playerHand");
  const tableCardsDiv = document.getElementById("tableCards");
  playerHandDiv.innerHTML = "";
  tableCardsDiv.innerHTML = "";

  // Display player's hand
  playerHand.forEach((card) => {
    const cardDiv = document.createElement("div");
    cardDiv.className = "card";
    cardDiv.innerHTML = `${card.rank} of ${card.suit}`;
    const throwBtn = document.createElement("button");
    throwBtn.className = "throw-card";
    throwBtn.innerText = "Throw Card";
    throwBtn.disabled = nextTurnId !== userId;
    throwBtn.addEventListener("click", function () {
      throwCard(card.id);
    });
    cardDiv.appendChild(throwBtn);
    playerHandDiv.appendChild(cardDiv);
  });

  // Display table cards
  tableCards.forEach((card) => {
    const cardDiv = document.createElement("div");
    cardDiv.className = "card";
    cardDiv.innerHTML = `${card.rank} of ${card.suit}`;
    tableCardsDiv.appendChild(cardDiv);
  });

  // Enable/Disable buttons based on the turn
  document.querySelectorAll(".throw-card").forEach((btn) => {
    btn.disabled = nextTurnId !== userId;
  });
}

function throwCard(cardId) {
  const throwCardMessage = {
    action: "throw_card",
    payload: { card_id: cardId },
    user_id: userId,
    room_id: roomId,
  };
  socket.send(JSON.stringify(throwCardMessage));
}
