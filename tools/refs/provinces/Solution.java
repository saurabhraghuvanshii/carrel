import java.util.*;

// Built from the dfs in the owner's Graphs/BreadthFirstSearch.java plus a counting loop. The
// owner's ConnectedComponenet.java does not compile: its dfsUtil calls a
// three-argument dfs that the file does not define, and its loop counts nothing.
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

    @SuppressWarnings("unchecked")
    private static ArrayList<Edge>[] build(int n, int[][] edges) {
        ArrayList<Edge>[] graph = new ArrayList[n];
        for (int i = 0; i < n; i++) {
            graph[i] = new ArrayList<>();
        }
        for (int[] e : edges) {
            graph[e[0]].add(new Edge(e[0], e[1], 1));
            graph[e[1]].add(new Edge(e[1], e[0], 1));
        }
        return graph;
    }

    public int countGroups(int n, int[][] edges) {
        ArrayList<Edge>[] graph = build(n, edges);
        boolean[] vis = new boolean[n];
        int groups = 0;
        for (int i = 0; i < n; i++) {
            if (!vis[i]) {
                groups++;
                dfs(graph, i, vis);
            }
        }
        return groups;
    }

    private void dfs(ArrayList<Edge>[] graph, int curr, boolean[] vis) {
        vis[curr] = true;
        for (int i = 0; i < graph[curr].size(); i++) {
            Edge e = graph[curr].get(i);
            if (!vis[e.dest]) {
                dfs(graph, e.dest, vis);
            }
        }
    }
}
