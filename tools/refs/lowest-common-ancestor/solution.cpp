static TreeNode* lca(TreeNode* n, int p, int q) {
    if (!n || n->val == p || n->val == q) return n;
    TreeNode* l = lca(n->left, p, q);
    TreeNode* r = lca(n->right, p, q);
    if (l && r) return n;
    return l ? l : r;
}

int commonAncestor(TreeNode* root, int p, int q) {
    return lca(root, p, q)->val;
}
