//import { NextApiRequest, NextApiResponse } from 'next';


export async function POST(req: Request) {
    console.log(req)
    
    const data = await fetch('https://api.vercel.app/blog')
    const posts = await data.json()

    return  new Response(JSON.stringify(posts) , { status: 200, })

  }

/*
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
    console.log(req)
    console.log(res)
    /*
    if (req.method !== 'POST') {
        return res.status(405).end(); // Method Not Allowed
    }
    * /
    return res.status(206).end();
    /*
    try {
        const { number, words } = req.body;

        // Process the data (replace this with your actual logic)
        const processedData = `Received number: ${number}, words: ${words.join(', ')}`; 

        res.status(200).json({ result: processedData });
    } catch (error) {
        console.error(error);
        res.status(500).json({ error: 'Internal Server Error' });
    }
    * /
}*/