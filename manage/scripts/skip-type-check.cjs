const fs = require('fs');
const path = require('path');

const generatedDir = './src/do-exercise-api'; // 生成的代码目录

function addTsNoCheck(dir) {
	fs.readdirSync(dir).forEach((file) => {
		const filePath = path.join(dir, file);
		const stat = fs.statSync(filePath);

		if (stat.isDirectory()) {
			addTsNoCheck(filePath); // 递归处理子目录
		} else if (filePath.endsWith('.ts')) {
			const content = fs.readFileSync(filePath, 'utf8');
			if (!content.startsWith('// @ts-nocheck')) {
				fs.writeFileSync(filePath, `// @ts-nocheck\n${content}`);
				console.log(`Added // @ts-nocheck to ${filePath}`);
			}
		}
	});
}

addTsNoCheck(generatedDir);