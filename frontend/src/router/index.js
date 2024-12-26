import Home from '../views/Home.vue';
import About from '../views/About.vue';
import Contact from '../views/Contact.vue';
import Redis from '../views/Redis.vue';
import MysqlUser from '../views/MysqlUser.vue';
import OfflineBinlog from '../views/OfflineBinlog.vue';
import { createRouter,createWebHistory } from 'vue-router';
const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/about', name: 'About', component: About },
  { path: '/contact', name: 'Contact', component: Contact },
  { path: '/redis', name: 'Redis', component: Redis },
  { path: '/mysqluser/:id', name: 'MysqlUser', component: MysqlUser },
  { path: '/offlinebinlog', name: 'OfflineBinlog', component: OfflineBinlog },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
  });
  
  export default router;
