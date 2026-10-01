#include <queue>
#include <utility>
#include <vector>
using namespace std;

int shortestPath(vector<vector<int>>& grid) {
    int rows = grid.size(), cols = grid[0].size();
    if (grid[0][0] || grid[rows - 1][cols - 1]) return -1;
    vector<vector<int>> dist(rows, vector<int>(cols));
    dist[0][0] = 1;
    queue<pair<int, int>> q;
    q.push({0, 0});
    int dy[] = {0, 1, 0, -1}, dx[] = {1, 0, -1, 0};
    while (!q.empty()) {
        auto [y, x] = q.front();
        q.pop();
        for (int d = 0; d < 4; d++) {
            int ny = y + dy[d], nx = x + dx[d];
            if (ny >= 0 && ny < rows && nx >= 0 && nx < cols && !grid[ny][nx] && !dist[ny][nx]) {
                dist[ny][nx] = dist[y][x] + 1;
                q.push({ny, nx});
            }
        }
    }
    return dist[rows - 1][cols - 1] ? dist[rows - 1][cols - 1] : -1;
}
