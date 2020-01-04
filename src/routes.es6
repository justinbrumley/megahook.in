import Home from './pages/Home';
import Inspect from './pages/Inspect';

const routes = [
  { path: '/', component: Home },
  { path: '/i/:name', component: Inspect },
];

export default routes;
