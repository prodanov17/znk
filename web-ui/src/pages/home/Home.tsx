import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

const Home = () => {
    const [userId, setUserId] = useState('');
    const [username, setUsername] = useState('');
    const [roomId, setRoomId] = useState('1234');

    useEffect(() => {
        const randomId = Math.floor(1000 + Math.random() * 9000);
        setUsername(`Player${randomId}`);
        setUserId(randomId.toString());
    }, []);

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
                <div className="flex justify-center">
                    <Link to={`/lobby/${roomId}?userId=${userId}&username=${username}`} className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-lg w-full">
                        Join Room
                    </Link>
                </div>
            </div>
        </div>
    );
}

export default Home;

