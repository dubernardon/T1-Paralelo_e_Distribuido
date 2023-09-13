#!/usr/bin/env python3
from random import uniform
from threading import Thread, Lock, Condition
from time import time, sleep

DIRTY = 0
CLEAN = 1

class ForkSyncManager:

    def __init__(self):
        self.condition = Condition()

    def wait(self):
        with self.condition:
            self.condition.wait()

    def notifyAll(self):
        with self.condition:
            self.condition.notify_all()

class Fork:

    def __init__(self, pid):

        self.pid = pid # philosopher's ID
        self.state = DIRTY
        self.lock = Lock()
        self.manager = ForkSyncManager()

    def request(self, pid):

        while self.pid != pid:
            if self.state == DIRTY:
                self.lock.acquire()
                self.state = CLEAN
                self.pid = pid
                self.lock.release()
            else:
                self.manager.wait()

    def done(self):

        self.state = DIRTY
        self.manager.notifyAll()

class Philosopher:

    def __init__(self, pid, left, right):

        self.pid = pid
        self.left = left
        self.right = right

        self.thread = Thread(target=self.dine)
        self.thread.start()

    def dine(self):
        self.think()
        self.eat()
        

    def eat(self):

        print(" "*self.pid*5 + f" R1" )
        self.left.request(self.pid)
        self.right.request(self.pid)

        
        sleep(3)
        print(" "*self.pid*5 + f" E1" )

        self.left.done()
        self.right.done()

    def think(self):
        
        sleep(2)
        print(" "*self.pid*5 + f" T1" )

class Table:

    def __init__(self):
        pass

    def start(self):
        self.forks = [Fork(p+1) for p in range(5)]
        self.philosophers = [Philosopher(i, self.forks[(i%5)],self.forks[(i+1)%5]) for i in range(5)]

if __name__ == "__main__":
    print("[P1] [P2] [P3] [P4] [P5]\n")
    table = Table()
    table.start()
