#include <queue>
#include <utility>
#include <vector>
using namespace std;

int minutesToInfectAll(vector<vector<int>>& grid) {
    int rows = grid.size(), cols = grid[0].size(), healthy = 0;
    queue<pair<int, int>> q;
    for (int r = 0; r < rows; r++)
        for (int c = 0; c < cols; c++) {
            if (grid[r][c] == 2) q.push({r, c});
            if (grid[r][c] == 1) healthy++;
        }
    int minutes = 0, dy[] = {0, 1, 0, -1}, dx[] = {1, 0, -1, 0};
    while (healthy > 0 && !q.empty()) {
        for (int size = q.size(); size > 0; size--) {
            auto [y, x] = q.front();
            q.pop();
            for (int d = 0; d < 4; d++) {
                int ny = y + dy[d], nx = x + dx[d];
                if (ny >= 0 && ny < rows && nx >= 0 && nx < cols && grid[ny][nx] == 1) {
                    grid[ny][nx] = 2;
                    healthy--;
                    q.push({ny, nx});
                }
            }
        }
        minutes++;
    }
    return healthy == 0 ? minutes : -1;
}
