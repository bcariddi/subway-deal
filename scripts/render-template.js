import fs from 'fs/promises';
import path from 'path';
import Handlebars from 'handlebars';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const templatesDir = path.join(__dirname, '..', 'templates');

// Register helpers
Handlebars.registerHelper('formatMoney', (value) => {
  if (value === 0) return '$0';
  const millions = value / 1000000;
  return `$${millions}M`;
});

Handlebars.registerHelper('eq', (a, b) => a === b);

Handlebars.registerHelper('lookup', (obj, index) => {
  if (Array.isArray(obj) && typeof index === 'number') {
    return obj[index];
  }
  return obj?.[index];
});

export async function renderTemplate(templateName, data) {
  const templatePath = path.join(templatesDir, `${templateName}.html`);
  const templateSource = await fs.readFile(templatePath, 'utf-8');
  const template = Handlebars.compile(templateSource);

  // Read base HTML
  const basePath = path.join(templatesDir, 'base.html');
  const baseHTML = await fs.readFile(basePath, 'utf-8');

  // Read CSS
  const cssPath = path.join(templatesDir, 'styles.css');
  const cssContent = await fs.readFile(cssPath, 'utf-8');

  // Render card content
  const cardHTML = template(data);

  // Inject CSS inline and card content into base
  const finalHTML = baseHTML
    .replace('<link rel="stylesheet" href="styles.css">', `<style>${cssContent}</style>`)
    .replace('<!-- Card content injected here -->', cardHTML);

  return finalHTML;
}
