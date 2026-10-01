#include <queue>
#include <vector>
using namespace std;

vector<vector<int>> levels(TreeNode* root) {
    vector<vector<int>> rows;
    queue<TreeNode*> q;
    if (root) q.push(root);
    while (!q.empty()) {
        vector<int> row;
        for (size_t i = q.size(); i > 0; i--) {
            TreeNode* n = q.front();
            q.pop();
            row.push_back(n->val);
            if (n->left) q.push(n->left);
            if (n->right) q.push(n->right);
        }
        rows.push_back(row);
    }
    return rows;
}
