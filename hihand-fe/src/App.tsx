import AppRouter from './router';
import './App.css';

const App = () => {
  return (
    <div>
      <nav style={{ padding: '1rem', borderBottom: '1px solid #ccc' }}>
        <a href="/" style={{ marginRight: 10 }}>Home</a>
      </nav>
      <AppRouter />
    </div>
  );
};

export default App;
