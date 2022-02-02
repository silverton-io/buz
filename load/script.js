import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('http://localhost:8080/i?stm=1643826582621&e=se&se_ca=ui&se_ac=clicked-button&se_la=struct-event&se_pr=someprop&se_va=118.99&eid=822405a1-381d-4a62-96b5-598684609cfe&tv=js-3.2.3&tna=sp&aid=site&p=web&cookie=1&cs=UTF-8&lang=en-US&res=1792x1120&cd=30&tz=America%2FNew_York&dtm=1643826582620&vp=1080x683&ds=1080x707&vid=1&sid=f0b4d620-5f14-4c02-be1d-75ebaa41b2a0&duid=2b153fcc-23fe-42d2-970d-a7e075bad47e&uid=jake%40bostata.com&url=http%3A%2F%2Flocalhost%3A8080%2F&cx=eyJzY2hlbWEiOiJpZ2x1OmNvbS5zbm93cGxvd2FuYWx5dGljcy5zbm93cGxvdy9jb250ZXh0cy9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6W3sic2NoZW1hIjoiaWdsdTpjb20uc25vd3Bsb3dhbmFseXRpY3Muc25vd3Bsb3cvd2ViX3BhZ2UvanNvbnNjaGVtYS8xLTAtMCIsImRhdGEiOnsiaWQiOiJkYzk4ZDA5Ny1mZWJiLTQ2M2ItYWUzMC0wNTAwMzBmODM0OWIifX1dfQ');
  sleep(1);
}