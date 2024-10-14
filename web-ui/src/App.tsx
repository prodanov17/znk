import { Route, Routes } from 'react-router-dom'
import Home from './pages/home/Home'
import Game from './pages/game/Game'
import Lobby from './pages/lobby/Lobby'
import { ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'

function App() {
    return (
        <>
            <ToastContainer
                position="top-right"
                autoClose={3000}
                hideProgressBar={false}
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                pauseOnHover
                theme="dark"
            />

            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/game/:room_id" element={<Game />} />
                <Route path="/lobby/:room_id" element={<Lobby />} />
                <Route path="*" element={<div>Not Found</div>} />
            </Routes>
        </>
    )
}

export default App
