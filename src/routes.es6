import Home from './pages/Home';
import Inspect from './pages/Inspect';

const routes = [
  { path: '/', component: Home },
  { path: '/m/:name/inspect', component: Inspect },
];

export default routes;