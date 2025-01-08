const fs = require('fs');
const path = require('path');

// Родительская директория, с которой начнём поиск (можно заменить на нужную)
const startDir = path.join(__dirname, '..');
// Итоговый файл, куда будет записываться информация
const outputFilePath = path.join(__dirname, 'go_files_output.txt');

// Функция рекурсивного обхода директорий
function walkDir(dir, callback) {
  fs.readdirSync(dir, { withFileTypes: true }).forEach((entry) => {
    const entryPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      walkDir(entryPath, callback);
    } else {
      callback(entryPath);
    }
  });
}

// Удаляем содержимое итогового файла, если он существует (чтобы начинать с чистого листа)
if (fs.existsSync(outputFilePath)) {
  fs.unlinkSync(outputFilePath);
}

// Запускаем обход директорий
walkDir(startDir, (filePath) => {
  // Проверяем расширение файла
  if (path.extname(filePath) === '.go') {
    try {
      const fileContent = fs.readFileSync(filePath, 'utf8');
      const fileName = path.relative(startDir, filePath);
      // Записываем в формат "имя_файла: текст"
      const lineToWrite = `=== ${fileName} ===\n${fileContent}\n\n`;
      fs.appendFileSync(outputFilePath, lineToWrite, 'utf8');
    } catch (err) {
      console.error(`Ошибка чтения файла ${filePath}: `, err);
    }
  }
});

console.log(`Сборка завершена. Результат записан в ${outputFilePath}`);
