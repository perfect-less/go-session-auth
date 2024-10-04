'use client';

import { FormEvent, useState } from "react";
import AcmeLogo from "../ui/acme-logo";
import { lusitana } from "../ui/fonts";
import { ArrowRightIcon, KeyIcon, UserIcon } from "@heroicons/react/24/outline";
import { Button } from "../ui/button";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const router = useRouter()
  const [isCredInvalid, setIsCredInvalid] = useState(false)

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const formData = new FormData(event.currentTarget)
    const response = await fetch('http://localhost:8082/login', {
      method: 'POST',
      credentials: 'include',
      body: formData,
    })

    if (response.ok) {
      console.log("OK")
      setIsCredInvalid(false)
      router.push("/dashboard")
    } else {
      console.log("NOT OK")
      setIsCredInvalid(true)
      // Handle errors
    }
  }

  return (
      <main className="flex items-center justify-center md:h-screen">
          <div className="relative mx-auto flex w-full max-w-[400px] flex-col space-y-2.5 p-4 md:-mt-32">
              <div className="flex h-20 w-full items-end rounded-lg bg-blue-500 p-3 md:h-36">
                  <div className="w-32 text-white md:w-36">
                      <AcmeLogo />
                  </div>
              </div>
              <form className="space-y-3" onSubmit={handleSubmit}>
                <div className="flex-1 rounded-lg bg-gray-50 px-6 pb-4 pt-8">
                  <h1 className={`${lusitana.className} mb-3 text-2xl`}>
                    Please log in to continue.
                  </h1>
                  {isCredInvalid == true && 
                    <div className="w-full">
                      <div className="mb-3 mt-5 py-3 px-3 bg-red-50 rounded-md border border-red-800 block text-xs font-medium text-gray-900">
                        <p className="text-red-900 text-xs">
                            Invalid username or password.
                        </p>
                      </div>
                    </div>
                  }
                  <div className="w-full">
                    <div>
                      <label
                        className="mb-3 mt-5 block text-xs font-medium text-gray-900"
                        htmlFor="username"
                      >
                        Username
                      </label>
                      <div className="relative">
                        <input
                          className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                          id="username"
                          type="username"
                          name="username"
                          placeholder="Enter your username"
                          required
                        />
                        <UserIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                      </div>
                    </div>
                    <div className="mt-4">
                      <label
                        className="mb-3 mt-5 block text-xs font-medium text-gray-900"
                        htmlFor="password"
                      >
                        Password
                      </label>
                      <div className="relative">
                        <input
                          className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm outline-2 placeholder:text-gray-500"
                          id="password"
                          type="password"
                          name="password"
                          placeholder="Enter password"
                          required
                          minLength={6}
                        />
                        <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                      </div>
                    </div>
                  </div>
                  <Button className="mt-4 w-full" type='submit'>
                    Log in <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
                  </Button>
                  <div className="flex h-8 items-end space-x-1">
                  {/* Add form errors here */}
                </div>
              </div>
            </form>
          </div>
      </main>
  );
}
