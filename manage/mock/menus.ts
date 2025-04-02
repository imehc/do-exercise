import { faker } from '@faker-js/faker';

export interface MenuItem {
  id: string;
  name: string;
  route: string;
  filePath: string;
  children?: MenuItem[];
}

/** 模拟获取管理员菜单 */
export const getAdminMenus = () => {
  return new Promise<MenuItem[]>((resolve) => {
    window.setTimeout(() => {
      const menus: MenuItem[] = Array.from({ length: 3 }, (_, index) => ({
        id: `p${index + 1}`,
        name: faker.word.sample(),
        route: `/page${index + 1}`,
        filePath: `/page${index + 1}/page.tsx`,
        children: Array.from({ length: 2 }, (_, idx) => ({
          id: `p${index + 1}c${idx + 1}`,
          name: faker.word.sample(),
          route: `/page${index + 1}/page${idx + 1}`,
          filePath: `/page${index + 1}/page${idx + 1}/page.tsx`,
        }))
      }));
      resolve(menus);
    }, 1000);
  });
};

/** 模拟获取普通菜单 */
export const getUserMenus = () => {
  return new Promise<MenuItem[]>((resolve) => {
    window.setTimeout(() => {
      const menus: MenuItem[] = Array.from({ length: 1 }, () => ({
        id: '1',
        name: faker.word.sample(),
        route: '/page1',
        filePath: `/page1/page.tsx`,
      }));
      resolve(menus);
    }, 1000);
  });
};