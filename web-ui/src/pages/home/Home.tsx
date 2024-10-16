import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';

const Home = () => {
    const [userId, setUserId] = useState('');
    const [username, setUsername] = useState('');
    const [roomId, setRoomId] = useState('1234');
    const navigate = useNavigate();

    useEffect(() => {
        const randomId = Math.floor(1000 + Math.random() * 9000);
        setUsername(`Player${randomId}`);
        setUserId(randomId.toString());
    }, []);

    const handleCreateRoom = async () => {
        try {
            const res = await fetch('http://localhost:8000/api/v1/ws/rooms', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    user_id: userId,
                    username: username,
                }),
            });

            if (res.ok) {
                const data = await res.json();
                setRoomId(data.room_id);
                toast.success('Room created successfully');
                navigate(`/lobby/${data.room_id}?userId=${userId}&username=${username}`);

            } else {
                toast.error('Failed to create room');
            }
        } catch (error) {
            toast.error('Failed to create room');
            console.error('Failed to create room:', error);
        }
    }

    const handleJoinRoom = async () => {
        try {
            if (!roomId) {
                toast.error('Room ID is required');
                return;
            }
            const res = await fetch(`http://localhost:8000/api/v1/ws/rooms/${roomId}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
            });
            console.log(res)
            if (res.ok) {
                console.log("HERE")
                const response = await res.json();
                console.log(res)
                if (response.clients.length >= 4) {
                    toast.error('Room is full');
                    return;
                }

                if (!response.room_id) {
                    toast.error('Room not found');
                    return;
                }
                navigate(`/lobby/${roomId}?userId=${userId}&username=${username}`);
            } else {
                toast.error('Room not found');
            }

        } catch (error) {
            toast.error('Failed to join room');
            console.error('Failed to join room:', error);
        }
    }


    return (
        <div className="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-center p-6">
            <div className="max-w-md w-full bg-gray-800 shadow-lg rounded-lg p-6 space-y-6">
                {/* Title */}
                <h1 className="text-3xl font-bold text-white text-center">Join the Game</h1>
                <p className="text-gray-400 text-center">Fill out the details below to enter a room.</p>

                {/* Username Input */}
                <input type="hidden" value={userId} />
                <div>
                    <label htmlFor="username" className="block text-sm text-gray-400 mb-2">Username</label>
                    <input
                        id="username"
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        className="w-full p-3 bg-gray-700 border border-gray-600 rounded text-gray-200 focus:ring-2 focus:ring-blue-600"
                    />
                </div>

                {/* Room ID Input */}
                <div>
                    <label htmlFor="room-id" className="block text-sm text-gray-400 mb-2">Room ID</label>
                    <input
                        id="room-id"
                        type="text"
                        value={roomId}
                        onChange={(e) => setRoomId(e.target.value)}
                        className="w-full p-3 bg-gray-700 border border-gray-600 rounded text-gray-200 focus:ring-2 focus:ring-blue-600"
                    />
                </div>

                {/* Button to Join Room */}
                <div className="flex flex-col items-center justify-center">
                    <button onClick={handleJoinRoom} className="bg-blue-600 text-center hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg w-full">
                        Join Room
                    </button>
                    <span className="flex items-center my-4 w-full">
                        <div className="h-px bg-gray-700 w-full"></div>
                        <span className="px-4 text-gray-400">OR</span>
                        <div className="h-px bg-gray-700 w-full"></div>
                    </span>

                    <button onClick={handleCreateRoom} className="bg-blue-600 text-center hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg w-full">
                        Create Room
                    </button>
                </div>
            </div>
        </div>
    );
}

export default Home;

