import Vue from 'vue'
import VueRouter from 'vue-router';
import routes from './routes.es6';

Vue.use(VueRouter);
Vue.config.productionTip = false

const router = new VueRouter({
  routes,
  mode: 'history',
});

new Vue({
  router,
}).$mount('#app')
