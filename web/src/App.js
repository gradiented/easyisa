import React from 'react';
import NavBar from './NavBar';
import { useAuth0 } from './Auth0Provider';

function App() {
  const { loading, getTokenSilently, user } = useAuth0();

  async function logToken() {
    if (!loading) {
      const token = await getTokenSilently();
      console.log('>>>>>>>>>>> ', token, ' <<<<<<<<<<<<');
      console.log('>>>>>>>>>>> ', user, '<<<<<<<<<<<<<');
    }
  }

  logToken();

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="App">
      <header>
        <NavBar />
      </header>
    </div>
  );
}

export default App;
