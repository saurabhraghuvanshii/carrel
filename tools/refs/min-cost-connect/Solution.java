import java.util.*;

// Ported from the owner's Graphs/PrimsAlgo.java (primsMst and Pair), unchanged; returns the
// cost it printed.
class Solution {
    static class Edge {
        int src;
        int dest;
        int wt;

        Edge(int s, int d, int w) {
            this.src = s;
            this.dest = d;
            this.wt = w;
        }
    }

    static class Pair implements Comparable<Pair> {
        int v;
        int cost;

        Pair(int v, int c) {
            this.v = v;
            this.cost = c;
        }

        @Override
        public int compareTo(Pair p2) {
            return this.cost - p2.cost;
        }
    }

    @SuppressWarnings("unchecked")
    private static ArrayList<Edge>[] build(int n, int[][] edges) {
        ArrayList<Edge>[] graph = new ArrayList[n];
        for (int i = 0; i < n; i++) {
            graph[i] = new ArrayList<>();
        }
        for (int[] e : edges) {
            graph[e[0]].add(new Edge(e[0], e[1], e[2]));
            graph[e[1]].add(new Edge(e[1], e[0], e[2]));
        }
        return graph;
    }

    public int cheapestConnection(int n, int[][] edges) {
        ArrayList<Edge>[] graph = build(n, edges);
        boolean[] vis = new boolean[graph.length];
        PriorityQueue<Pair> pq = new PriorityQueue<>();
        pq.add(new Pair(0, 0));
        int finalCost = 0;
        while (!pq.isEmpty()) {
            Pair curr = pq.remove();
            if (!vis[curr.v]) {
                vis[curr.v] = true;
                finalCost += curr.cost;
                for (int i = 0; i < graph[curr.v].size(); i++) {
                    Edge e = graph[curr.v].get(i);
                    pq.add(new Pair(e.dest, e.wt));
                }
            }
        }
        return finalCost;
    }
}
