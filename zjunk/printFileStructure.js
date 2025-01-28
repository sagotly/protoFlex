// Импортируем необходимые модули
const fs = require('fs');
const path = require('path');

// Функция для построения структуры дерева
function getDirectoryTree(dir, prefix = '') {
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  entries.forEach((entry, index) => {
    const isLast = index === entries.length - 1;
    const entryPath = path.join(dir, entry.name);
    const connector = isLast ? '└── ' : '├── ';

    console.log(`${prefix}${connector}${entry.name}`);

    if (entry.isDirectory()) {
      const newPrefix = prefix + (isLast ? '    ' : '│   ');
      getDirectoryTree(entryPath, newPrefix);
    }
  });
}

// Указываем начальную директорию
const startDirectory = path.resolve(__dirname);

// Выводим файловую структуру проекта
console.log(startDirectory);
getDirectoryTree(startDirectory);