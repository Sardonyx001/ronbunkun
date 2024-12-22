"use client"

import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";

export default function HomePage() {
  const [maxPapers, setMaxPaperCount] = useState('');
  const [keywords, setKeywords] = useState('');
  const [result, setResult] = useState('');

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    //const keywordArray = keywords.split(/\s+/); 

    console.log(JSON.stringify({ maxPapers, keywords })) 

    try {
      const response = await fetch('http://localhost:3001/api/generate', { 
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ keywords, maxPapers }) 
      });

      const data = await response.json();
      setResult(data.result); 
    } catch (error) {
      setResult('An error occurred.'); 
      console.error(error);
    }
  };

  return (
    <main className="container mx-auto p-4"> 
      <h1 className="text-3xl font-bold">Input Page</h1>

      <form className="flex-col" onSubmit={handleSubmit}>
        <label 
          className='my-4'htmlFor="maxPapers">Max Page Count:</label>
        <Input 
          className='my-4'
          type="number" 
          id="maxPapers" 
          value={maxPapers} 
          onChange={(e) => setMaxPaperCount(e.target.value)} 
        />

        <label htmlFor="keywords">Keywords:</label>
        <Textarea 
          className='my-4'
          id="keywords" 
          value={keywords} 
          onChange={(e) => setKeywords(e.target.value)} 
        />

        <Button type="submit">Submit</Button>
      </form>

      {result && (
        <div className="mt-4">
          <h2>Result:</h2>
          <p>{result}</p> 
        </div>
      )}
    </main>
  );
}