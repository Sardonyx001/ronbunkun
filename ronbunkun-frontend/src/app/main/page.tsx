"use client"

import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

type Article = {
  id: string,
  title: string,
  published: string,
  pdfurl: string,
  categories: string[],
  summary: string
}

export default function HomePage() {
  const [maxPapers, setMaxPaperCount] = useState('');
  const [keywords, setKeywords] = useState('');
  const [result, setResult] = useState <Article[]>();
  interface GenerateData {
    keywords: string[];
    maxPapers: number;
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const keywordArray = keywords.split(/\s+/); 

    //console.log(JSON.stringify({ maxPapers, keywords })) 

    const genData : GenerateData = {
        keywords: keywordArray,
        maxPapers: parseInt(maxPapers)
    }

    console.log(JSON.stringify(genData)) 

    try {
      const response = await fetch('http://127.0.0.1:3001/api/generate', { 
        method: 'POST',
        //mode: 'no-cors',
        headers: { 'Content-Type': 'application/json'
        //,'Access-Control-Allow-Origin': '*'
      },
        body: JSON.stringify(genData) 
      });

      const data = await response.json();
      setResult(data.result); 
      console.log(data)
    } catch (error) {
      //setResult(data.error); 
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

      {result && 
        (<Table>
          <TableCaption>A list of your recent invoices.</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">Invoice</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Method</TableHead>
              <TableHead className="text-right">Amount</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {result.map((item) => (
              <TableRow key={item.id}>
                <TableCell className="font-medium">{item.pdfurl}</TableCell>
                <TableCell>{item.title}</TableCell>
                <TableCell>{item.published}</TableCell>
                <TableCell className="">{item.summary}</TableCell>
              </TableRow>
            ))}
          </TableBody>
          <TableFooter>
            <TableRow>
              <TableCell colSpan={3}>Total</TableCell>
              <TableCell className="text-right">$2,500.00</TableCell>
            </TableRow>
          </TableFooter>
        </Table>)
        }
    </main>
  );
}