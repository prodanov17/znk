import { useParams, useLocation, useNavigate } from "react-router-dom";
import { ClipboardIcon, SwitchIcon } from "../../components/Icons";
import { useEffect, useRef, useState } from "react";
import { WebSocketService } from "../../services/websocketService";
import { toast } from "react-toastify";

export type Team = {
    team_id: number;
    players: Player[];
    score: number;
};

export type Player = {
    user_id: string;
    username: string;
    score: number;
};


type Message = {
    username: string;
    message: string;
}

const Lobby = () => {
    const { room_id } = useParams();
    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const userId = queryParams.get("userId");
    const username = queryParams.get("username");

    const navigate = useNavigate();

    const wsServiceRef = useRef<WebSocketService | null>(null);


    // State to store player count and teams data
    const [playerCount, setPlayerCount] = useState(0);
    const [teams, setTeams] = useState<Team[]>([]);
    const [messages, setMessages] = useState<Message[]>([]);
    const [messageBox, setMessageBox] = useState("");

    const [rules, setRules] = useState<{ [key: string]: string }>({});

    // Handle rule changes dynamically
    const handleRuleChange = (e: React.ChangeEvent<HTMLInputElement>, ruleKey: string) => {
        setRules((prevRules) => ({
            ...prevRules,
            [ruleKey]: e.target.value, // Update specific rule based on ruleKey
        }));
    };


    useEffect(() => {
        console.log("Lobby component mounted!", room_id, userId, username);
        if (room_id && userId && username) {
            console.log(`Connecting to WebSocket for room ${room_id}...`);
            const wsService = WebSocketService.getInstance(`/api/v1/ws/join/${room_id}?userId=${userId}&username=${username}`);
            wsService.connect();

            // Handle player_joined action
            wsService.on("player_joined", (message) => {
                const { playerCount, teams } = message.payload;
                setPlayerCount(parseInt(playerCount)); // Update player count
                setTeams(teams);
                console.log("setting teams", teams)
                toast.info("A player joined the game!");
            });

            wsService.on("player_left", (message) => {
                const { playerCount, teams } = message.payload;
                setPlayerCount(parseInt(playerCount)); // Update player count
                setTeams(teams); // Update teams
                toast.warn("A player left the game!");
            });



            wsService.on("team_changed", (message) => {
                const { teams } = message.payload;
                setTeams(teams);
            }
            );

            wsService.on("teams", (message) => {
                const { teams } = message.payload;
                setTeams(teams);
            })

            wsService.on("message", (message) => {
                const { username, message: msg } = message.payload;
                setMessages((prevMessages) => [...prevMessages, { username, message: msg }]);
            });

            wsService.on("game_started", (message) => {
                const { dealer_id } = message.payload;
                navigate(`/game/${room_id}?userId=${userId}&username=${username}&dealer_id=${dealer_id}`);
            });

            // Store the service instance in the ref
            wsServiceRef.current = wsService;

            console.log(teams)
            setTimeout(() => sendActiontoWS("get_teams", {}), 500);
            // Cleanup WebSocket on unmount
            return () => {
                console.log("lobby component unmounted!");
                if (wsServiceRef.current) {
                    // wsServiceRef.current.close()
                }
            };

        }
    }, [room_id, userId, username]);


    const sendMessagetoWS = () => {
        const payload = {
            action: "message",
            payload: {
                message: messageBox
            },
            user_id: userId,
            room_id: room_id
        }
        if (wsServiceRef.current) {
            wsServiceRef.current.sendMessage(payload);
        }

        setMessageBox("");
    }

    const sendActiontoWS = (action: string, payload: Record<string, any>) => {
        const packetPayload = {
            action: action,
            payload: payload,
            user_id: userId,
            room_id: room_id
        }
        if (wsServiceRef.current) {
            wsServiceRef.current.sendMessage(packetPayload);
        }
    }


    return (
        <div className="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-between p-6">
            {/* Header Section */}
            <header className="w-full text-center mb-6">
                <h1 className="text-3xl font-bold text-white">Game Lobby</h1>
                <h2 className="text-lg font-semibold text-gray-400 space-x-2">
                    <span>Join Code: {room_id}</span>
                    <button className="" onClick={() => {
                        navigator.clipboard.writeText(room_id || "")
                        toast.info("Join code copied to clipboard!")
                    }}>
                        <ClipboardIcon />
                    </button>
                </h2>
                {playerCount < 4 ? <p className="text-gray-400">Waiting for players...</p> : <p className="text-green-400">Ready to start!</p>}
                <p className="text-gray-300 font-medium">Players: {playerCount}/4</p> {/* Display player count */}
            </header>

            {/* Teams Section */}
            <div className="flex w-full max-w-3xl justify-between mb-6">
                {/* Team 1 */}
                <div className="flex flex-col w-1/2 bg-gray-800 shadow-md rounded-lg p-4 mx-2">
                    <h3 className="text-lg font-semibold text-gray-300 mb-4">Team 1</h3>
                    <div className="flex flex-col space-y-2">
                        {/* Render players in Team 1 */}
                        {teams[0]?.players?.length > 0 ? (
                            teams[0].players.map((player, index) => (
                                <p key={index} className={`rounded p-2 text-gray-200 ${player.user_id == userId ? "bg-blue-600" : "bg-gray-700"}`}>
                                    {player.username}
                                </p>
                            ))
                        ) : (
                            <p className="text-gray-400">No players yet...</p>
                        )}
                    </div>
                </div>

                {/* Switch Teams Button */}
                <div className="flex items-center justify-center mx-4">
                    <button className="bg-blue-600 hover:bg-blue-700 text-white rounded-full p-3" onClick={() => sendActiontoWS("change_team", {})}>
                        <SwitchIcon className="w-6 h-6" />
                    </button>
                </div>

                {/* Team 2 */}
                <div className="flex flex-col w-1/2 bg-gray-800 shadow-md rounded-lg p-4 mx-2">
                    <h3 className="text-lg font-semibold text-gray-300 mb-4">Team 2</h3>
                    <div className="flex flex-col space-y-2">
                        {/* Render players in Team 2 */}
                        {teams[1]?.players?.length > 0 ? (
                            teams[1].players.map((player, index) => (
                                <p key={index} className={`rounded p-2 text-gray-200 ${player.user_id == userId ? "bg-blue-600" : "bg-gray-700"}`}>
                                    {player.username}
                                </p>
                            ))
                        ) : (
                            <p className="text-gray-400">No players yet...</p>
                        )}
                    </div>
                </div>
            </div>
            <div className="w-full max-w-xl bg-gray-800 shadow-md rounded-lg p-6 mb-6">
                <h2 className="text-lg font-semibold text-gray-300 mb-4">Game Rules</h2>
                <label htmlFor="rules" className="block text-sm text-gray-400 mb-2">Play until:</label>
                <input
                    id="rules"
                    type="text"
                    className="w-full p-2 bg-gray-700 border border-gray-600 rounded text-gray-200 focus:ring-2 focus:ring-blue-600"
                    placeholder="e.g., 10 points"
                    onChange={(e) => handleRuleChange(e, 'winning_score')}
                    value={rules['winning_score'] || ''}
                />
            </div>


            {/* Chat Section */}
            <div className="w-full max-w-3xl bg-gray-800 shadow-md rounded-lg p-6 flex flex-col space-y-4">
                <div className="h-32 overflow-y-auto bg-gray-700 p-4 rounded">
                    {messages.length > 0 ? messages.map((msg, index) => (
                        <p key={index} className="text-gray-200">{msg.username}: {msg.message}</p>
                    )) : (<p className="text-gray-400"> Start chatting...</p>)
                    }
                </div>
                <div className="flex space-x-2">
                    <input
                        type="text"
                        className="w-full p-2 bg-gray-700 border border-gray-600 rounded text-gray-200 focus:ring-2 focus:ring-blue-600"
                        placeholder="Type a message..."
                        onChange={(e) => setMessageBox(e.target.value)}
                        value={messageBox}
                    />
                    <button className="bg-blue-600 hover:bg-blue-700 text-white p-2 rounded" onClick={sendMessagetoWS} >
                        Send
                    </button>
                </div>
            </div>

            {/* Start Game Button */}
            <div className="w-full max-w-3xl flex justify-end mt-6">
                <button onClick={() => sendActiontoWS("start_game", { rules: { winning_score: rules['winning_score'] } })} className="bg-green-600 hover:bg-green-700 text-white font-semibold p-3 rounded-lg">
                    Start Game
                </button>
            </div>
        </div>
    );
};

export default Lobby;

