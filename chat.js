import OpenAI from "openai";
import { promises as fs } from 'fs';
import path from 'path';

const token = process.env["GITHUB_TOKEN"];
const endpoint = "https://models.inference.ai.azure.com";
const modelName = "o3-mini";

async function readFiles() {
  const files = ['context/openapi.yaml', 'golang/sdk-example.go', 'context/changelog.json'];
  const fileContents = {};

  for (const file of files) {
    try {
      const content = await fs.readFile(path.join(process.cwd(), file), 'utf8');
      fileContents[file] = content;
    } catch (error) {
      console.error(`Error reading ${file}:`, error);
      fileContents[file] = `Error: Could not read ${file}`;
    }
  }
  return fileContents;
}

export async function main() {
  const fileContents = await readFiles();
  
  const enhancedPrompt = `Context: openapi.yaml file represents MongoDB Atlas OpenAPI. 
sdk-example.go is the file I want to edit. 
changelog.json contains new changes to the attached OpenAPI.yaml

SDK uses OpenAPI tag and operationId fields to build SDK methods.

File contents:

openapi.yaml:
${fileContents['openapi.yaml']}

sdk-example.go:
${fileContents['sdk-example.go']}

changelog.json:
${fileContents['changelog.json']}

Action: Act as software agent that for provided SDK example extending it by providing sdk method calls from changelog.json path`;

  const client = new OpenAI({ baseURL: endpoint, apiKey: token });

  const response = await client.chat.completions.create({
    messages: [
      { role: "developer", content: "You are a helpful assistant." },
      { role: "user", content: enhancedPrompt }
    ],
    model: modelName
  });

  console.log(response.choices[0].message.content);
}

main().catch((err) => {
  console.error("The sample encountered an error:", err);
});
