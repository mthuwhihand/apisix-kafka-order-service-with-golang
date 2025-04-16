import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from '../pages/Home/Home';

const AppRouter = () => (
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<Home />} />
        </Routes>
    </BrowserRouter>
);

export default AppRouter;
