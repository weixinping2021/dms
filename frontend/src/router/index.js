import DbList from '../views/DbList.vue';
import About from '../views/About.vue';
import Contact from '../views/Contact.vue';
import Redis from '../views/Redis.vue';
import Mysql from '../views/Mysql.vue';
import OfflineBinlog from '../views/OfflineBinlog.vue';
import { createRouter,createWebHistory } from 'vue-router';
const routes = [
  { path: '/', name: 'DbList', component: DbList },
  { path: '/about', name: 'About', component: About },
  { path: '/contact', name: 'Contact', component: Contact },
  { path: '/redis', name: 'Redis', component: Redis },
  { path: '/mysql/:id', name: 'Mysql', component: Mysql },
  { path: '/offlinebinlog', name: 'OfflineBinlog', component: OfflineBinlog },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
  });
  
  export default router;
