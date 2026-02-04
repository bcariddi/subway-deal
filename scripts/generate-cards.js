import puppeteer from 'puppeteer';
import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import { renderTemplate } from './render-template.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.join(__dirname, '..');

const OUTPUT_DIR = path.join(projectRoot, 'output', 'cards');
const DATA_DIR = path.join(projectRoot, 'data');

// Card dimensions with bleed at 300 DPI
// 2.5" + 3mm bleed on each side = 2.5" + 0.236" = 2.736"
// 3.5" + 3mm bleed on each side = 3.5" + 0.236" = 3.736"
const CARD_WIDTH_PX = 820;  // ~2.73" at 300 DPI
const CARD_HEIGHT_PX = 1120; // ~3.73" at 300 DPI

async function generateCard(browser, cardData, templateName, outputFilename) {
  const page = await browser.newPage();

  // Set viewport to exact card dimensions
  await page.setViewport({
    width: CARD_WIDTH_PX,
    height: CARD_HEIGHT_PX,
    deviceScaleFactor: 1,
  });

  // Render template
  const html = await renderTemplate(templateName, cardData);
  await page.setContent(html, { waitUntil: 'networkidle0' });

  // Screenshot
  const outputPath = path.join(OUTPUT_DIR, outputFilename);
  await page.screenshot({
    path: outputPath,
    type: 'png',
    omitBackground: false,
  });

  await page.close();
  console.log(`  Generated: ${outputFilename}`);
}

async function expandCards(cards) {
  // Expand cards with quantity > 1
  const expanded = [];
  for (const card of cards) {
    const quantity = card.quantity || 1;
    for (let i = 0; i < quantity; i++) {
      expanded.push({
        ...card,
        instanceNumber: i + 1,
      });
    }
  }
  return expanded;
}

async function loadJSON(filename) {
  const filePath = path.join(DATA_DIR, filename);
  const content = await fs.readFile(filePath, 'utf-8');
  return JSON.parse(content);
}

async function main() {
  console.log('Subway Deal Card Generator');
  console.log('==========================\n');

  // Ensure output directory exists
  await fs.mkdir(OUTPUT_DIR, { recursive: true });

  const browser = await puppeteer.launch({ headless: true });

  try {
    // Load all data files
    const properties = await loadJSON('properties.json');
    const wildcards = await loadJSON('wildcards.json');
    const actions = await loadJSON('actions.json');
    const rent = await loadJSON('rent.json');
    const money = await loadJSON('money.json');

    let cardCount = 0;

    // Generate property cards
    console.log('Generating property cards...');
    for (const card of properties.cards) {
      await generateCard(
        browser,
        card,
        'property-card',
        `${card.id}.png`
      );
      cardCount++;
    }

    // Generate wildcard cards
    console.log('\nGenerating wildcard cards...');
    for (const card of wildcards.cards) {
      await generateCard(
        browser,
        card,
        'wildcard-card',
        `${card.id}.png`
      );
      cardCount++;
    }

    // Generate action cards (expand by quantity)
    console.log('\nGenerating action cards...');
    const expandedActions = await expandCards(actions.cards);
    for (const card of expandedActions) {
      await generateCard(
        browser,
        card,
        'action-card',
        `${card.id}_${card.instanceNumber}.png`
      );
      cardCount++;
    }

    // Generate rent cards (expand by quantity)
    console.log('\nGenerating rent cards...');
    const expandedRent = await expandCards(rent.cards);
    for (const card of expandedRent) {
      await generateCard(
        browser,
        card,
        'rent-card',
        `${card.id}_${card.instanceNumber}.png`
      );
      cardCount++;
    }

    // Generate money cards (expand by quantity)
    console.log('\nGenerating money cards...');
    const expandedMoney = await expandCards(money.cards);
    for (const card of expandedMoney) {
      await generateCard(
        browser,
        card,
        'money-card',
        `${card.id}_${card.instanceNumber}.png`
      );
      cardCount++;
    }

    console.log(`\n${'='.repeat(40)}`);
    console.log(`Generated ${cardCount} cards successfully`);

    if (cardCount !== 106) {
      console.warn(`Warning: Expected 106 cards, generated ${cardCount}`);
    } else {
      console.log('All 106 cards generated!');
    }

  } finally {
    await browser.close();
  }
}

main().catch(console.error);
