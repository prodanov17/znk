import { useEffect, useRef, useState } from "react";
import Card from "../../components/Card";
import { WebSocketService } from "../../services/websocketService";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Player } from "../lobby/Lobby";
import { toast } from "react-toastify";
import GameRulesHover from "../../components/GameRulesHover";

type Card = {
    id: number;
    suit: string;
    rank: string;
};

type GameState = {
    next_turn_id?: string;
    dealer_id?: string;
};

type Team = {
    team_id: number;
    players: Player[];
    score: number;
};


const Game = () => {
    const navigate = useNavigate();

    const { room_id } = useParams();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const userId = queryParams.get("userId");
    const username = queryParams.get("username");

    const [playerHand, setPlayerHand] = useState<Card[]>([]);
    const [tableCards, setTableCards] = useState<Card[]>([]);
    const [tableValue, setTableValue] = useState<number>(0);
    const [dreamCard, setDreamCard] = useState<Card | null>(null);
    const [gameInfo, setGameInfo] = useState<Record<string, string>>({});

    const [teams, setTeams] = useState<Team[]>([]);

    const [gameState, setGameState] = useState<GameState>({});
    const [playing, setPlaying] = useState(false);

    const [playerHandCount, setPlayerHandCount] = useState<{ [userId: string]: number }>({});
    const wsServiceRef = useRef<WebSocketService | null>(null);


    useEffect(() => {
        const wsService = WebSocketService.getInstance(`/api/v1/ws/join/${room_id}?userId=${userId}&username=${username}`);

        wsService.connect();

        // Set up event handlers for messages
        wsService.on("message", (message) => {
            console.log("New message received:", message.payload);
            // Handle the new message
        });

        wsService.on("deal_cards", (message) => {
            console.log("Dealing cards:", message.payload);
            // Update player's hand with the new cards
            setPlayerHand(message.payload.cards);
            setTableCards(message.payload.table_cards);
            setTableValue(message.payload.table_value);
            setGameState((prev) => ({ ...prev, next_turn_id: message.payload.next_turn_id }));
            setPlayerHandCount(message.payload.players_card_count);
            setDreamCard(message.payload.dream_card);
            setPlaying(true)
        })

        wsService.on("card_played", (message) => {
            if (message.payload.take_cards) {
                toast.info(`${message.payload.username} took with ${message.payload.card.rank} of ${message.payload.card.suit}`)
            }
            if (message.payload.value >= 10) {
                toast.info("ZNK! 🚀")
            }

            setTableCards(message.payload.table_cards);
            setTableValue(message.payload.table_value);
            setPlayerHandCount(message.payload.players_card_count);
            setPlaying(message.payload.playing)
            setPlayerHand(message.payload.player_hand);
            setGameState((prev) => ({ ...prev, next_turn_id: message.payload.next_turn_id }));

            if (message.payload.next_turn_id === userId) {
                toast.info("Your turn!")
            }
        })

        wsService.on("round_over", (message) => {
            const { dealer_id } = message.payload

            setGameState((prev) => ({ ...prev, dealer_id: dealer_id }))
            setTableCards([]);
            setTableValue(0);
            setTeams(message.payload.teams);
            setPlaying(false);

            toast.info(`Round over! Cards go to team ${message.payload.last_capture_id}`)
        })

        wsService.on("game_ended", (message) => {
            console.log("Game ended:", message.payload);
            const { teams, winner_team } = message.payload;

            setTeams(teams);

            const winnerTeam = teams.find((team: Team) => team.team_id === winner_team);
            toast.info(`Team ${winnerTeam?.team_id} (${winnerTeam.players[0].username} and ${winnerTeam.players[1].username}) wins with ${winnerTeam.score}!`);

            setTimeout(() => {
                navigate(`/lobby/${room_id}?userId=${userId}&username=${username}`);
            }, 5000);
        })

        wsService.on("game_state", (message) => {
            const payload = message.payload;
            setGameState({ dealer_id: payload.dealer_id, next_turn_id: payload.next_turn_id });
            setPlayerHandCount(payload.players_card_count);
            setPlayerHand(payload.player_hand);
            setTeams(payload.teams);
            setPlaying(!payload.playing);
            setTableCards(payload.table_cards);
            setTableValue(payload.table_value);
            setDreamCard(payload.dream_card);
            setGameInfo(payload.game_info);

            if (payload.state === "waiting") {
                setPlaying(false);
                navigate(`/lobby/${room_id}?userId=${userId}&username=${username}`);
            }

        })

        wsService.on("teams", (message) => {
            const { teams } = message.payload;

            setTeams(teams);

        })

        wsServiceRef.current = wsService;
        setTimeout(() => sendActiontoWS("game_state"), 500)

        // Clean up on component unmount
        return () => {
            console.log("Game component unmounted!");
            if (wsServiceRef.current) {
                // wsServiceRef.current.close()
            }
        };
    }, [room_id, userId, username]);


    const getPlayerTeamIndex = (playerId: string) => {
        return teams.findIndex((team) => team.players.some((player) => player.user_id === playerId));
    }

    const sendActiontoWS = (action: string, payload?: Record<string, any>) => {
        console.log("Sending action to WS:", action);
        const sendPayload = {
            action: action,
            payload: { payload },
            user_id: userId,
            room_id: room_id
        }
        if (wsServiceRef.current) {
            wsServiceRef.current.sendMessage(sendPayload);
        }
    }

    const getTeammateIndex = (playerId: string) => {
        const teamIndex = getPlayerTeamIndex(playerId);
        return teams[teamIndex].players.findIndex((player) => player.user_id !== playerId);
    }

    const getTeammate = (playerId: string) => {
        return teams[getPlayerTeamIndex(playerId)].players[getTeammateIndex(playerId)];
    }


    const getOpponentTeamIndex = (playerId: string) => {
        return getPlayerTeamIndex(playerId) === 0 ? 1 : 0;
    }

    const getTeam = (playerId: string) => {
        return teams[getPlayerTeamIndex(playerId)];
    }


    if (wsServiceRef.current === null) {
        return <div>Connecting to the server...</div>;
    }

    const teammate = getTeammate(userId || "");
    const myTeam = getTeam(userId || "");
    console.log(playerHandCount)
    return (
        <>
            <GameRulesHover rules={gameInfo} />
            {dreamCard &&
                <div className="flex justify-center items-center space-x-1 flex-col max-w-[300px] mx-auto pt-4 absolute right-0">
                    <Card suit={dreamCard?.suit || "back"} rank={dreamCard?.rank || "0"} className="w-20 " />
                    <Card
                        suit="back"
                        rank="0"
                        className="w-20 -mt-24 shadow"
                    />
                </div>}
            <section className="flex justify-center flex-col items-center w-screen overflow-hidden">
                <div className="flex justify-center items-center space-x-1 flex-col max-w-[300px] mx-auto pt-4">
                    <div className="flex justify-center items-center space-x-1">
                        {playerHandCount[teammate.user_id] > 0 ? Array(playerHandCount[teammate.user_id] || 4).fill(0).map((_, index) => (
                            <Card
                                key={index}
                                suit="back"
                                rank="0"
                                className="w-20"
                                style={{
                                    marginLeft: `-${Math.min(index, 1) * 60}px`,
                                }}

                            />
                        )) : <Card suit="back" rank="0" className="w-20 opacity-0" />}
                    </div>
                    <div className="flex flex-col text-center mt-4">
                        <p className="text-xs uppercase font-light text-neutral-400"> Team {myTeam.team_id} ({myTeam.score})</p>
                        <p className={`font-semibold flex items-center gap-2 transition-colors ${teammate.user_id == gameState.next_turn_id ? "text-white" : "text-neutral-500"}`}> {myTeam.players[getTeammateIndex(userId || "")].username} {gameState.dealer_id == teammate.user_id && "(D)"} </p>

                    </div>
                </div>
                <div className="flex h-[400px] justify-between items-center w-full px-4">
                    {teams[getOpponentTeamIndex(userId || "")].players.map((player, index) => (
                        <div key={index} className="flex basi-1/2 justify-center items-center flex-col  max-w-[200px] -translae-x-1/4 mx-aut">
                            {playerHandCount[player.user_id] > 0 ? Array(playerHandCount[player.user_id] || 4).fill(0).map((_, index) => (
                                <Card key={index} suit="back" rank="0" className="w-20 -mt-24" />
                            )) : Array(4).fill(0).map((_, index) => (<Card key={index} suit="back" rank="0" className="w-20 -mt-24 opacity-0" />))}
                            <div className="flex flex-col text-center mt-4">
                                <p className="text-xs uppercase font-light text-neutral-400"> Team {teams[getOpponentTeamIndex(userId || "")].team_id} ({teams[getOpponentTeamIndex(userId || "")].score})</p>
                                <p className={`font-semibold flex items-center gap-2 text-neutral-400 transition-colors ${player.user_id == gameState.next_turn_id ? "text-white" : "text-neutral-500"}`}> {player.username} {gameState.dealer_id == player.user_id && "(D)"}</p>
                            </div>
                        </div>
                    ))}
                </div>
                <div className="flex justify-center items-center space-x-1 flex-col max-w-[300px] mx-auto">
                    <div className="flex justify-center items-center space-x-1">
                        {playerHand ? playerHand.map((card, index) => (
                            <button key={index} onClick={() => sendActiontoWS("throw_card", { "card_id": card.id })} className="w-20">

                                <Card
                                    suit={card.suit}
                                    rank={card.rank}
                                    className="w-20"
                                />
                            </button>)) : <Card suit="back" rank="0" className="w-20 opacity-0" />}
                    </div>
                    <div className="flex flex-col text-center items-center mt-4">
                        <p className="text-xs uppercase font-light text-neutral-400"> Team {getTeam(userId || "").team_id} ({getTeam(userId || "").score})</p>
                        <p className={`font-semibold flex items-center gap-2 text-neutral-400 transition-colors ${userId === gameState.next_turn_id ? "text-white" : "text-neutral-500"}`}> You {gameState.dealer_id == userId && "(D)"}</p>
                    </div>

                </div>


            </section>
            <section className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 flex justify-center flex-col w-screen overflow-hidden">
                <div className="flex justify-center items-center space-x-1">
                    {tableCards && tableCards.slice(-4).map((card, index) => (
                        <Card
                            key={index}
                            suit={card.suit}
                            rank={card.rank}
                            className=" w-20"
                            style={{
                                marginLeft: `-${Math.min(index, 1) * 50}px`, // if index less than length - 4, add huge margin
                            }}
                        />
                    ))}
                </div>
                <div className="flex flex-col text-center mt-4">
                    <p className="text-xs uppercase font-light text-neutral-400"> {tableCards.length} cards</p>
                    <p className="text-sm">Value: <span className="font-semibold">{tableValue}</span></p>
                </div>


            </section>


            {gameState && gameState.dealer_id === userId && !playing && <button onClick={() => sendActiontoWS("deal_cards")} className="fixed bottom-4 right-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg">Deal Cards</button>}
        </>

    );
}

export default Game;
