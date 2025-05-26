const fs = require('fs');
const path = require('path');

const generatedDir = './src/do-exercise-api'; // 生成的代码目录

function processFile(filePath) {
	const content = fs.readFileSync(filePath, 'utf8');
	let updatedContent = content;
	let modified = false;

	// 1. 添加 // @ts-nocheck（如果没有）
	if (!updatedContent.startsWith('// @ts-nocheck')) {
		updatedContent = `// @ts-nocheck\n${updatedContent}`;
		modified = true;
		console.log(`🛡️  Added // @ts-nocheck to ${filePath}`);
	}

	// 2. 查找未导入的 ToJSON/FromJSON 函数
	const matches = [...updatedContent.matchAll(/(\b\w+(ToJSON|FromJSON)\b)/g)];
	const usedFns = new Set(matches.map((m) => m[1]));

	// 已存在的导入
	const existingImports = new Set(
		[...updatedContent.matchAll(/import\s+{([^}]+)}\s+from\s+['"]([^'"]+)['"]/g)]
			.flatMap((m) => m[1].split(',').map((s) => s.trim()))
	);

	const importLines = [];

	usedFns.forEach((fnName) => {
		if (!existingImports.has(fnName)) {
			const baseName = fnName.replace(/(ToJSON|FromJSON)$/, '');
			importLines.push(`import { ${fnName} } from './${baseName}';`);
		}
	});

	if (importLines.length > 0) {
		const lines = updatedContent.split('\n');
		lines.splice(1, 0, ...importLines); // 插入到第二行（保留 @ts-nocheck）
		updatedContent = lines.join('\n');
		modified = true;
		console.log(`🔧 Fixed missing imports in ${filePath}`);
	}

	if (modified) {
		fs.writeFileSync(filePath, updatedContent, 'utf8');
	}
}

function processDirectory(dir) {
	fs.readdirSync(dir).forEach((file) => {
		const filePath = path.join(dir, file);
		const stat = fs.statSync(filePath);

		if (stat.isDirectory()) {
			processDirectory(filePath);
		} else if (filePath.endsWith('.ts')) {
			processFile(filePath);
		}
	});
}

processDirectory(generatedDir);
