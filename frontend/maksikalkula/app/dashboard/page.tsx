'use client';

import React from "react";
import { Button } from "../ui/button";
import { useRouter } from "next/navigation";



export default function Page() {
  const router = useRouter()
  const [session_status, setSession_status] = React.useState('Unkown')

  async function handleClick () {
    const response = await fetch('http://localhost:8082/session_check', {
      method: 'GET',
      credentials: 'include'
    })
    if (response.ok) {
      setSession_status("Valid")
    } else {
      setSession_status("Invalid")
      router.push("/login")
    }
  }

  return (
    <main className="flex items-center justify-center md:h-screen">
      <p>
         Session: {session_status}
      </p>
      <br />
      <Button onClick={handleClick}>
         Click here
      </Button>
    </main>
  )
}
